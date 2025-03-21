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

const httpErrorStatusCodeWithMessage = {
  400: "BadRequest",
  401: "Unauthorized",
  404: "NotFound",
  500: "InternalServerError",
};
document.addEventListener("htmx:afterRequest", (event) => {
  const httpResponseObject = event.detail.xhr;

  if (
    httpResponseObject.getResponseHeader("Content-Type") !== "application/json"
  ) {
    return;
  }

  const responseData = httpResponseObject.responseText;
  if (responseData === "") {
    return;
  }

  const parsedResponse = JSON.parse(responseData);
  if (parsedResponse.body === undefined || parsedResponse.body === "") {
    return;
  }

  httpResponseStatusCode = httpResponseObject.status;

  let toastType = "success";
  let toastMessage = "Success";
  if (httpResponseStatusCode == 207) {
    toastType = "partialSuccess";
    toastMessage = "PartialSuccess";
  }

  if (httpResponseStatusCode >= 400) {
    toastType = "danger";
    toastMessage = httpErrorStatusCodeWithMessage[httpResponseStatusCode];
  }

  if (typeof parsedResponse.body === "string") {
    toastMessage = parsedResponse.body;
  }

  Alpine.store("toast").displayToast(toastMessage, toastType);
});
