import combobox from './combobox.js';

export default () => ({
    quack: new Set([0, 1, 2, 3]),
    pathways: new Set(),

    table: {
        ['x-if']() { return this.pathways.size > 0 }
    },

    rows: {
        ['x-for']: "item in Array.from(pathways).sort()",
        [':key']: "item.id"
    },

    row: {
        ['x-text']: "item.toString()",
    },

    remove: {
        ['@click']() {
            this.pathways.forEach((item) => {
                if (item.id === this.$el.dataset.id) {
                    this.pathways.delete(item)
                }
            })
        },
        [':data-id']: "item.id"
    },

    ...combobox({
        placeholder: "Add Pathway…",

        init: function () {
            Alpine.effect(() => {
                if (Alpine.store('species') != null) {
                    this.data().then(data => this.results = data)
                    this.pathways.clear()
                }
            })
        },

        callback: function (value) {
            this.pathways.add(value)
            this.value = ''
            if (this.pathways.size > 0) {
                var values = this.pathways.values()
                var next = values.next()
            }
        },

        data: async function (query) {
            var species = this.$store.species
            if (species == null || species.id == null) {
                return
            }

            var uri = `/api/v1/kegg/pathway?o=${species.id}`
            if (query != null && query !== "") {
                uri = `${uri}&q=${query}`
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
                            return item.id + " - " + item.name;
                        }
                        return item
                    });
                }
            }).catch(error => { console.log(error); });
        }
    })
})
