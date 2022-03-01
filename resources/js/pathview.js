'use strict';

import Alpine from 'alpinejs';
import Papa from 'papaparse';
import file from './file.js'
import pathway from './pathway.js';
import species from './species.js';
import '../css/main.css';

window.Alpine = Alpine;
window.Papa = Papa;

// Alpine.data('combobox', () => ({
//     ...combobox({
//         data: ["Species 1", "Species 2", "Species 3"]
//     })
// }));

Alpine.data('file', file)
Alpine.data('pathway', pathway);
Alpine.data('species', species);

Alpine.start();

