export default () => ({
    name: '',

    init() { },

    get have_file() {
        return this.name != '';
    },

    chip: {
        ['x-text']() { return this.$refs.file.files[0].name; }
    },

    remove: {
        ['@click']() {
            this.$refs.file.value = ''
            this.name = ''
        }
    },

    input: {
        ['x-ref']: "file",
        ['@change']() {
            if (this.$refs.file.files.length > 0) {
                this.name = this.$refs.file.files[0].name;
            } else {
                this.name = ''
            }
        }
    },

    options: {
        ['x-cloak']: "",
        ['x-transition']: "",
        [':class']: "have_file ? '' : 'hidden'"
    }
})
