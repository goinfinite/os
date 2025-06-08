UiToolset.RegisterAlpineState(() => {
  Alpine.data("marketplace", () => ({
    // PrimaryState
    marketplaceItem: {},
    get hostnameWithTrailingSlash() {
      return this.marketplaceItem.hostname + "/";
    },
    get dataFieldsAsString() {
      let dataFieldsAsString = "";
      for (let dataField of this.marketplaceItem.dataFields) {
        if (!dataField.value) {
          continue;
        }

        dataFieldsAsString += dataField.name + ":" + dataField.value + ";";
      }
      return dataFieldsAsString.slice(0, -1);
    },
    resetPrimaryStates() {
      this.marketplaceItem = {
        id: "",
        name: "",
        hostname: "",
        urlPath: "",
        dataFields: [],
      };
    },
    init() {
      this.resetPrimaryStates();
    },

    // AuxiliaryState
    selectedMarketplaceCatalogVerticalTab: "apps",
    updateSelectedMarketplaceCatalogVerticalTab(tabName) {
      this.selectedMarketplaceCatalogVerticalTab = tabName;
    },
    reloadMarketplacePageContent(listType) {
      htmx.ajax("GET", "/marketplace/?listType=" + listType, {
        select: "#marketplace-page-content",
        target: "#marketplace-page-content",
        indicator: "#loading-overlay",
        swap: "outerHTML transition:true",
      });
    },
    imageLightbox: {
      isOpen: false,
      imageUrl: "",
    },
    openImageLightbox(imageUrl) {
      this.imageLightbox.isOpen = true;
      this.imageLightbox.imageUrl = imageUrl;
    },
    closeImageLightbox() {
      this.imageLightbox.isOpen = false;
      this.imageLightbox.imageUrl = "";
    },

    // ModalState
    isScheduleSelectedMarketplaceItemInstallationModalOpen: false,
    openScheduleSelectedMarketplaceItemInstallationModal(catalogItemId) {
      this.resetPrimaryStates();

      const catalogItemEntity = JSON.parse(
        document.getElementById("marketplaceCatalogItem_" + catalogItemId)
          .textContent
      );
      this.marketplaceItem.id = catalogItemId;
      this.marketplaceItem.name = catalogItemEntity.name;

      for (const dataField of catalogItemEntity.dataFields) {
        dataField.value = dataField.defaultValue;
        this.marketplaceItem.dataFields.push(dataField);
      }

      this.isScheduleSelectedMarketplaceItemInstallationModalOpen = true;
    },
    closeScheduleSelectedMarketplaceItemInstallationModal() {
      this.isScheduleSelectedMarketplaceItemInstallationModalOpen = false;
    },
    isUninstallMarketplaceItemModalOpen: false,
    openUninstallMarketplaceItemModal(installedItemId, installedItemName) {
      this.resetPrimaryStates();

      this.marketplaceItem.id = installedItemId;
      this.marketplaceItem.name = installedItemName;
      this.isUninstallMarketplaceItemModalOpen = true;
    },
    closeUninstallMarketplaceItemModal() {
      this.isUninstallMarketplaceItemModalOpen = false;
    },
    uninstallMarketplaceItem() {
      htmx
        .ajax(
          "DELETE",
          "/api/v1/marketplace/installed/" + this.marketplaceItem.id + "/",
          { swap: "none" }
        )
        .then(() => {
          this.$store.main.refreshScheduledTasksPopover();
          this.$dispatch("uninstall:marketplace-item");
        })
        .finally(() => {
          this.closeUninstallMarketplaceItemModal();
        });
    },
  }));
});
