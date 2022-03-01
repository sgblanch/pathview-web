// import Alpine from 'alpinejs';
import combobox from './combobox.js';

export default () => ({
    ...combobox({
        placeholder: "Select Species…",

        callback: function (value) {
            Alpine.store('species', value);
        },

        data: async function (query) {
            var uri = `/api/v1/kegg/organism`
            if (query != null && query !== "") {
                uri = `${uri}?q=${query}`
            }

            return fetch(uri).then(response => {
                if (!response.ok) {
                    throw new Error("HTTP error, status = " + response.status);
                }
                return response.json()
            }).then(data => {
                if (data != null) {
                    return data.map(function (item) {
                        item.toString = function () {
                            if (item.common != null) {
                                return item.code + " - " + item.name + " - " + item.common;
                            }
                            return item.code + " - " + item.name;
                        }
                        return item
                    });
                }
            }).catch(error => { console.log(error); });
        }
    })
})
