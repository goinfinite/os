document.addEventListener("alpine:init", () => {
  Alpine.data("mappingSecurityRules", () => ({
    mappingSecurityRule: {},
    resetPrimaryStates() {
      this.mappingSecurityRule = {
        id: "",
        name: "",
        description: "",
        allowedIpsText: "",
        blockedIpsText: "",
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
    prepareFormData() {
      if (this.mappingSecurityRule.allowedIpsText) {
        const allowedIps = this.mappingSecurityRule.allowedIpsText
          .split("\n")
          .map((ip) => ip.trim())
          .filter((ip) => ip !== "");

        if (allowedIps.length > 0) {
          document
            .querySelector("form")
            .insertAdjacentHTML(
              "beforeend",
              `<input type="hidden" name="allowedIps" value='${JSON.stringify(
                allowedIps
              )}' />`
            );
        }
      }

      if (this.mappingSecurityRule.blockedIpsText) {
        const blockedIps = this.mappingSecurityRule.blockedIpsText
          .split("\n")
          .map((ip) => ip.trim())
          .filter((ip) => ip !== "");

        if (blockedIps.length > 0) {
          document
            .querySelector("form")
            .insertAdjacentHTML(
              "beforeend",
              `<input type="hidden" name="blockedIps" value='${JSON.stringify(
                blockedIps
              )}' />`
            );
        }
      }
    },

    isCreateMappingSecurityRuleModalOpen: false,
    openCreateMappingSecurityRuleModal() {
      this.resetPrimaryStates();
      this.isCreateMappingSecurityRuleModalOpen = true;
    },
    closeCreateMappingSecurityRuleModal() {
      this.prepareFormData();
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

      this.mappingSecurityRule.id = mappingSecurityRuleId;
      this.mappingSecurityRule.name = mappingSecurityRuleEntity.name;

      if (mappingSecurityRuleEntity.description) {
        this.mappingSecurityRule.description =
          mappingSecurityRuleEntity.description;
      }

      if (
        mappingSecurityRuleEntity.allowedIps &&
        mappingSecurityRuleEntity.allowedIps.length > 0
      ) {
        this.mappingSecurityRule.allowedIpsText =
          mappingSecurityRuleEntity.allowedIps.join("\n");
      }

      if (
        mappingSecurityRuleEntity.blockedIps &&
        mappingSecurityRuleEntity.blockedIps.length > 0
      ) {
        this.mappingSecurityRule.blockedIpsText =
          mappingSecurityRuleEntity.blockedIps.join("\n");
      }

      if (mappingSecurityRuleEntity.rpsSoftLimitPerIp) {
        this.mappingSecurityRule.rpsSoftLimitPerIp =
          mappingSecurityRuleEntity.rpsSoftLimitPerIp;
      }

      if (mappingSecurityRuleEntity.rpsHardLimitPerIp) {
        this.mappingSecurityRule.rpsHardLimitPerIp =
          mappingSecurityRuleEntity.rpsHardLimitPerIp;
      }

      if (mappingSecurityRuleEntity.responseCodeOnMaxRequests) {
        this.mappingSecurityRule.responseCodeOnMaxRequests =
          mappingSecurityRuleEntity.responseCodeOnMaxRequests;
      }

      if (mappingSecurityRuleEntity.maxConnectionsPerIp) {
        this.mappingSecurityRule.maxConnectionsPerIp =
          mappingSecurityRuleEntity.maxConnectionsPerIp;
      }

      if (mappingSecurityRuleEntity.responseCodeOnMaxConnections) {
        this.mappingSecurityRule.responseCodeOnMaxConnections =
          mappingSecurityRuleEntity.responseCodeOnMaxConnections;
      }

      if (mappingSecurityRuleEntity.bandwidthBpsLimitPerConnection) {
        this.mappingSecurityRule.bandwidthBpsLimitPerConnection =
          mappingSecurityRuleEntity.bandwidthBpsLimitPerConnection;
      }

      if (mappingSecurityRuleEntity.bandwidthLimitOnlyAfterBytes) {
        this.mappingSecurityRule.bandwidthLimitOnlyAfterBytes =
          mappingSecurityRuleEntity.bandwidthLimitOnlyAfterBytes;
      }

      this.isUpdateMappingSecurityRuleModalOpen = true;
    },
    closeUpdateMappingSecurityRuleModal() {
      this.prepareFormData();
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
