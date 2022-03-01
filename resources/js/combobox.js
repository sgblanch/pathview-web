export default (config) => ({
    _open: false,
    focused: 0,
    value: '',
    prev: '',
    placeholder: config.placeholder ?? "Select...",
    callback: config.callback,
    child_init: config.init,
    results: [],

    data: config.data,

    init() {
        switch (typeof this.data) {
            case 'function':
                this.$watch('value', ((value) => {
                    this.data(value).then(data => this.results = data)
                }))
                this.data().then(data => this.results = data)
                break;
            case 'object':
                this.results = config.data
                this.$watch('value', ((value) => {
                    this.filter(value)
                }))
                break;
            default:
                console.log("unhandled type:", typeof this.data)
        }

        if (typeof this.child_init == 'function') {
            this.child_init()
        }
    },

    filter(needle) {
        if (!this._open || !needle) return this.results = this.data

        this.results = this.data.filter((value) => value.toLowerCase().includes(needle.toLowerCase()))
        this.focused = 0
    },

    open() {
        this.prev = this.value
        this.$el.select()
        this._open = true
    },

    cancel() {
        if (!this._open) {
            return
        }

        this.$el.blur()
        this._open = false
        this.value = this.prev
    },

    commit() {
        if (this.results == null || this.results.length == 0) {
            this.cancel()
            return
        }

        this.$el.blur()
        this._open = false
        this.value = this.results[this.focused].toString();

        if (typeof this.callback == 'function') {
            this.callback(this.results[this.focused])
        }
    },

    combobox: {
        ['@click.outside']() { this.cancel(); },
    },

    search: {
        [':placeholder']: "placeholder",
        ['@click']() {
            if (this._open) {
                this.commit()
            } else {
                this.open()
            }
        },
        ['@keydown.escape']() { this.cancel(); },
        ['@keydown.arrow-up.prevent']() { this.focused = Math.max(0, this.focused - 1); },
        ['@keydown.arrow-down.prevent']() { this.focused = Math.min(this.results.length - 1, this.focused + 1); },
        ['@keydown.enter']() { this.commit(); },
        ['x-model']: "value",
        ['x-ref']: "search",
    },

    options: {
        [':hidden']: "!_open",
    },

    option: {
        [':class']: "{ 'text-white bg-indigo-600': index === focused, 'text-gray-900': index !== focused, 'icon-check': value === item }",
        ['x-text']: "item.toString()",
        ['@click']() { this.commit(); },
        ['@mouseover']: "focused = index;",
        [':id']: "item.id",
    }
})
