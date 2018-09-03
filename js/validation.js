function loadValidation() {
    $.validator.addMethod("match", function(value, element, params) {
        return this.optional(element) || $("#"+params["target"]).val() == value
    }, $.validator.format("Must match {target}"));
}
