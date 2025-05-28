Infinite.RegisterAlpineState(dropzoneAlpineState);

function dropzoneAlpineState() {
  Alpine.data("dropzone", () => ({
    files: [],
    updateFileInput() {
      const dataTransfer = new DataTransfer();
      this.files.forEach((file) => dataTransfer.items.add(file));
      this.$refs.dropzone.files = dataTransfer.files;
    },
    removeFile(index) {
      this.files.splice(index, 1);
      this.updateFileInput();
    },
    handleDrop(event) {
      const droppedFiles = Array.from(event.dataTransfer.files);
      this.files = this.files.concat(droppedFiles);
      this.updateFileInput();
    },
    resetState() {
      this.files = [];
      this.updateFileInput();
    },
    init() {
      this.resetState();

      document.addEventListener("delete:dropzone-state", () => {
        this.resetState();
      });
    },
  }));
}
