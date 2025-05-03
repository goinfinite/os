document.addEventListener("alpine:init", () => {
  Alpine.data("mappingSecurityRules", () => ({
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

    get shouldDisableCreateMappingSecurityRuleSubmitButton() {
      return this.mappingSecurityRule.name === "";
    },
    get shouldDisableUpdateMappingSecurityRuleSubmitButton() {
      return this.mappingSecurityRule.name === "";
    },

    isCreateMappingSecurityRuleModalOpen: false,
    openCreateMappingSecurityRuleModal() {
      this.resetPrimaryStates();
      this.isCreateMappingSecurityRuleModalOpen = true;
    },
    closeCreateMappingSecurityRuleModal() {
      this.isCreateMappingSecurityRuleModalOpen = false;
    },

    isUpdateMappingSecurityRuleModalOpen: false,
    openUpdateMappingSecurityRuleModal(mappingSecurityRuleId) {
      this.resetPrimaryStates();

      const mappingSecurityRuleEntity = JSON.parse(
        document.getElementById(
          "mappingSecurityRuleEntity_" + mappingSecurityRuleId
        ).textContent
      );
      this.mappingSecurityRule = mappingSecurityRuleEntity;
      this.isUpdateMappingSecurityRuleModalOpen = true;
    },
    closeUpdateMappingSecurityRuleModal() {
      this.isUpdateMappingSecurityRuleModalOpen = false;
    },

    isDeleteMappingSecurityRuleModalOpen: false,
    openDeleteMappingSecurityRuleModal(
      mappingSecurityRuleId,
      mappingSecurityRuleName
    ) {
      this.resetPrimaryStates();
      this.mappingSecurityRule.id = mappingSecurityRuleId;
      this.mappingSecurityRule.name = mappingSecurityRuleName;
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
});
