const path = require('path');
const fs = require('fs');
const svgstore = require('svgstore');
const heroicons = path.dirname(require.resolve('heroicons/README.md'));

var sprites = svgstore()
    .add('adjustments', fs.readFileSync(path.resolve(heroicons, 'outline/adjustments.svg'), 'utf8'))
    .add('check', fs.readFileSync(path.resolve(heroicons, 'outline/check.svg'), 'utf8'))
    // .add('document-add', fs.readFileSync(path.resolve(heroicons, 'solid/document-add.svg'), 'utf8'))
    .add('document-text', fs.readFileSync(path.resolve(heroicons, 'solid/document-text.svg'), 'utf8'))
    .add('plus', fs.readFileSync(path.resolve(heroicons, 'solid/plus.svg'), 'utf8'))
    // .add('upload', fs.readFileSync(path.resolve(heroicons, 'outline/upload.svg'), 'utf8'))
    .add('x', fs.readFileSync(path.resolve(heroicons, 'solid/x.svg'), 'utf8'));

fs.writeFileSync('../static/sprites.svg', sprites);
