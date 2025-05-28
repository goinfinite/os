Infinite.RegisterAlpineState(mappingSecurityRulesIndexAlpineState);

function mappingSecurityRulesIndexAlpineState() {
  Alpine.data("mappingSecurityRules", () => ({
    // PrimaryState
    mappingSecurityRule: {},
    resetPrimaryStates() {
      this.mappingSecurityRule = {
        id: "",
        name: "",
        description: "",
        allowedIps: [],
        blockedIps: [],
        rpsSoftLimitPerIp: "",
        rpsHardLimitPerIp: "",
        responseCodeOnMaxRequests: "429",
        maxConnectionsPerIp: "",
        bandwidthBpsLimitPerConnection: "",
        bandwidthLimitOnlyAfterBytes: "",
        responseCodeOnMaxConnections: "420",
      };
    },
    init() {
      this.resetPrimaryStates();
    },

    // ModalState
    isCreateMappingSecurityRuleModalOpen: false,
    openCreateMappingSecurityRuleModal() {
      this.resetPrimaryStates();
      this.isCreateMappingSecurityRuleModalOpen = true;
    },
    closeCreateMappingSecurityRuleModal() {
      this.isCreateMappingSecurityRuleModalOpen = false;
    },

    isUpdateMappingSecurityRuleModalOpen: false,
    openUpdateMappingSecurityRuleModal(secRuleId) {
      this.resetPrimaryStates();

      const secRuleEntity = JSON.parse(
        document.getElementById("secRuleEntity_" + secRuleId).textContent
      );
      this.mappingSecurityRule = secRuleEntity;
      this.isUpdateMappingSecurityRuleModalOpen = true;
    },
    closeUpdateMappingSecurityRuleModal() {
      this.isUpdateMappingSecurityRuleModalOpen = false;
    },

    isDeleteMappingSecurityRuleModalOpen: false,
    openDeleteMappingSecurityRuleModal(secRuleId, secRuleName) {
      this.resetPrimaryStates();
      this.mappingSecurityRule.id = secRuleId;
      this.mappingSecurityRule.name = secRuleName;
      this.isDeleteMappingSecurityRuleModalOpen = true;
    },
    closeDeleteMappingSecurityRuleModal() {
      this.isDeleteMappingSecurityRuleModalOpen = false;
    },
    deleteMappingSecurityRule() {
      this.closeDeleteMappingSecurityRuleModal();
      htmx
        .ajax(
          "DELETE",
          "/api/v1/vhost/mapping/security-rule/" +
            this.mappingSecurityRule.id +
            "/",
          {
            swap: "none",
          }
        )
        .then(() => this.$dispatch("refresh:mapping-security-rules-table"));
    },
  }));
}
