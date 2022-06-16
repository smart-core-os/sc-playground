import {execSync} from 'child_process';
import replace from 'replace-in-file';

const protoFiles = [
  'pkg/device/evcharger/evcharger.proto',
  'pkg/playpb/playground.proto'
];
const protocPluginOpts = '--js_out=import_style=commonjs:. --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:.';

const out = execSync(`protomod protoc -- -I../.. ${protocPluginOpts} ${protoFiles.join(' ')}`);
console.log(out.toString());

// update the generated files to replace
// `require('../../../traits/*_pb.js');`
// with `require('@smart-core-os/sc-api-grpc-web/...')`

// replace .js imports
replace.sync({
  files: ['pkg/**/*_pb.js', 'trait/**/*_pb.js'],
  from: /require\('(?:\.\.\/){3,}((?:traits|types|info)\/.+_pb.js)'\)/g,
  to: `require('@smart-core-os/sc-api-grpc-web/$1')`
});
// replace .d.ts imports
replace.sync({
  files: ['pkg/**/*_pb.d.ts', 'trait/**/*_pb.d.ts'],
  from: /from '(?:\.\.\/){3,}((?:traits|types|info)\/.+_pb)'/g,
  to: `from '@smart-core-os/sc-api-grpc-web/$1'`
});
