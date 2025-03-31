document.addEventListener("alpine:init", () => {
  Alpine.data("ssls", () => ({
    // Primary states
    sslPair: {},
    resetPrimaryStates() {
      this.sslPair = {
        id: "",
        virtualHostsHostnames: [],
        certificate: "",
        key: "",
      };
    },
    init() {
      this.resetPrimaryStates();
    },

    // Auxiliary states
    get shouldDisableImportSslCertificateSubmitButton() {
      return (
        this.sslPair.virtualHostsHostnames.length == 0 ||
        this.sslPair.certificate == "" ||
        this.sslPair.key == ""
      );
    },
    shouldImportSslCertificateAsFile: false,
    get shouldDisableRemoveVirtualHostsHostnamesSubmitButton() {
      return (
        this.sslPair.id == "" || this.sslPair.virtualHostsHostnames.length == 0
      );
    },
    downloadPemFile(isKeyFileContent) {
      let pemFileContent = this.sslPair.certificate;
      let fileExtension = "crt";
      if (isKeyFileContent) {
        pemFileContent = this.sslPair.key;
        fileExtension = "key";
      }

      const pemFileNameWithExtension = `${this.sslPair.id}.${fileExtension}`;
      Infinite.DownloadFile(pemFileNameWithExtension, pemFileContent);
    },

    // Modal states
    isImportSslCertificateModalOpen: false,
    openImportSslCertificateModal() {
      this.resetPrimaryStates();

      this.isImportSslCertificateModalOpen = true;
    },
    closeImportSslCertificateModal() {
      this.sslPair.certificate = btoa(this.sslPair.certificate);
      this.sslPair.key = btoa(this.sslPair.key);

      this.isImportSslCertificateModalOpen = false;
    },
    isViewPemFilesModalOpen: false,
    openViewPemFilesModal(sslPairId, certificateContent, keyContent) {
      this.resetPrimaryStates();

      this.sslPair.id = sslPairId;
      this.sslPair.certificate = certificateContent;
      this.sslPair.key = keyContent;
      this.isViewPemFilesModalOpen = true;
    },
    closeViewPemFilesModal() {
      this.isViewPemFilesModalOpen = false;
    },
    isSwapToSelfSignedModalOpen: false,
    openSwapToSelfSignedModal(sslPairId) {
      this.resetPrimaryStates();

      this.sslPair.id = sslPairId;
      this.isSwapToSelfSignedModalOpen = true;
    },
    closeSwapToSelfSignedModal() {
      this.isSwapToSelfSignedModalOpen = false;
    },
    swapToSelfSigned() {
      htmx
        .ajax("DELETE", "/api/v1/ssl/" + this.sslPair.id + "/", {
          swap: "none",
        })
        .finally(() => this.closeSwapToSelfSignedModal());
    },
  }));
});
