import toml from 'toml'
import fs from 'fs'
import path from 'path'
import { spawnSync } from 'child_process'
import { Storage } from '@google-cloud/storage';

import {
  crTemplatesPath,
  templatesBucketName,
} from './secrets.js'

const storage = new Storage()

const dirPath = path.join('templates')

const isDev = process.env.NODE_ENV === 'dev' || process.argv[2] === 'dev'

function getAllTemplateDirectories() {
  return fs
    .readdirSync(dirPath)
    .map(p => {
      const templateDirPath = path.join(dirPath, p)
      const templateDirStats = fs.statSync(templateDirPath)
      if (templateDirStats.isDirectory()) return templateDirPath
    })
    .filter(p => !!p)
}

function getTemplate(directory) {
  console.log('Parsing template.toml -', directory)

  const templatePath = path.join(directory, 'template.toml')
  try {
    return toml.parse(fs.readFileSync(templatePath))
  } catch (error) {
    throw new Error(`Error parsing template.toml on path "${templatePath}"`, error)
  }
}

function buildTemplate(directory, template) {
  console.log('Building template image -', template.id)

  const code_cells_path = template.code_cells_dir.split(path.sep)

  if (code_cells_path[0] !== template.files_dir) {
    throw new Error(`Parent directory of code_cells_dir ("${code_cells_dir}") must be a "${template.files_dir}"`)
  }

  const name = `${crTemplatesPath}/${template.id}`
  const image = isDev ? `${name}:dev` : `${name}:latest`

  const root_dir = '/home/runner'

  const repo_files_dir = path.join(directory, template.files_dir)

  const files_dir = '/home/setup'
  const code_cells_dir = path.join(root_dir, ...code_cells_path.slice(1))

  const buildArgs = Object.entries({
    root_dir,
    repo_files_dir,
    files_dir,
    setup_cmd: template.setup_cmd,
    start_cmd: template.start_cmd,
    code_cells_dir,
  })
    .flatMap(e => ['--build-arg', `"${e[0]}"="${e[1]}"`])

  const tagArgs = ['-t', image]

  const command = `docker build .`

  const result = spawnSync(
    command,
    [...buildArgs, ...tagArgs],
    { stdio: 'inherit', shell: true },
  )

  if (result.status !== 0) {
    throw new Error(`Error building template image "${template.id}"`)
  }

  return {
    id: template.id,
    image,
    root_dir,
    code_cells_dir,
  }
}

function pushTemplate(template) {
  console.log('Pushing template image -', template.image)

  const command = `docker push "${template.image}"`
  const result = spawnSync(command, { stdio: 'inherit', shell: true })

  if (result.status !== 0) {
    throw new Error(`Error pushing template image "${template.image}"`)
  }

  return template
}

async function uploadTemplate(template) {
  if (isDev) {
    console.log('NOT uploading the `template.toml` config file -', template.id, 'only the `npm run build` uploads template configs.')
    return
  }
  console.log('Uploading template config -', template.id)
  try {
    const file = storage.bucket(templatesBucketName).file(`${template.id}.json`)
    await file.save(JSON.stringify(template), {
      metadata: {
        contentType: 'application/json',
        cacheControl: 'no-cache, max-age=0',
      },
      public: true,
      validation: 'md5',
    })
  } catch (error) {
    throw new Error(`Error uploading template "${template.id}" config`, error)
  }
}

try {
  console.log('Publishing templates...')
  await Promise.all(getAllTemplateDirectories()
    .map(templateDir => ({ template: getTemplate(templateDir), directory: templateDir }))
    .map(({ template, directory }) => buildTemplate(directory, template))
    .map(template => pushTemplate(template))
    .map(template => uploadTemplate(template))
  )
} catch (error) {
  console.error(error)
  process.exitCode = 1
}
