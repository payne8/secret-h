const proxy = require('http-proxy-middleware');
const Bundler = require('parcel-bundler');
const express = require('express');
const { spawn } = require('child_process');

let bundler = new Bundler('./index.html', {
  hmr: false
});
let app = express();
let secretHServer;

function startSecretHServer() {
  console.log('starting secret-h');
  secretHServer = spawn('../myapp', []);
  secretHServer.stdout.on('data', out => console.log(out.toString()));
  secretHServer.stderr.on('data', out => console.error(out.toString()));
  secretHServer.on('close', console.log);
}

startSecretHServer();

app.use(
  '/api',
  proxy({
    logLevel: 'info',
    target: 'http://localhost:8080'
  })
);

app.get('/reset', (req, res, next) => {
  secretHServer.kill();
  startSecretHServer();
  setTimeout(() => {
    res.status(200).send('Done');
  }, 3000);
});

app.use(bundler.middleware());

app.listen(Number(process.env.PORT || 1234));

process.on('exit', () => {
  secretHServer.kill();
});
