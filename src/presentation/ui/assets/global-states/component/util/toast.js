document.addEventListener("alpine:init", () => {
  Alpine.store("toast", {
    toastVisible: false,
    toastMessage: "",
    toastType: "danger",

    displayToast(message, type) {
      this.toastVisible = true;
      this.toastMessage = message;
      this.toastType = type;
      setTimeout(() => {
        this.clearToast();
      }, 3000);
    },

    clearToast() {
      this.toastVisible = false;
      this.toastMessage = "";
    },
  });
});

document.addEventListener("htmx:afterRequest", (event) => {
  const httpResponseObject = event.detail.xhr;

  const contentType = httpResponseObject.getResponseHeader("Content-Type");
  if (contentType !== "application/json") {
    return;
  }

  const responseData = httpResponseObject.responseText;
  if (responseData === "") {
    return;
  }

  let toastType = "success";
  const isResponseError = httpResponseObject.status >= 400;
  if (isResponseError) {
    toastType = "danger";
  }

  if (httpResponseObject.status == 207) {
    console.log(httpResponseObject);
  }

  const parsedResponse = JSON.parse(responseData);
  if (parsedResponse.body === undefined || parsedResponse.body === "") {
    return;
  }
  const toastMessage = parsedResponse.body;

  Alpine.store("toast").displayToast(toastMessage, toastType);
});
