import path from 'node:path';

export default {
  schemaPath: path.resolve(process.cwd(), 'config/openapi.json'),
  serversPath: path.resolve(process.cwd(), 'src/services'),
  projectName: 'clear-bill',
  namespace: 'API',
  enumStyle: 'string-literal',
  declareType: 'interface',
  splitDeclare: false,
  requestImportStatement: "import { request } from '../request';",
  requestOptionsType: '{ [key: string]: any }',
  dataFields: ['data'],
  isCamelCase: true,
};
