
function send_ajax(data, url, callback) {
	$.ajax({
		type: "POST",
        url: url,
        data: data,
		success: callback,
		error: function(error) {
            console.log("Error doing ajax request...");
            console.log(error);
            $(".errmsg").text(error.responseJSON["error"]);
            $("#error").removeClass("invisible");
		}
	});
}

function add_secret() {
    var secret = $("[name=secret]").val();
    var password = $("[name=password]").val();
    var expire = $("[name=expire]").val();
    var max_view = $("[name=maxview]").val();
    var url = window.location.href;
    data = {"secret": secret, "password": password, "expire": expire, "maxview": max_view}
    send_ajax(data, "/secret", function(result) {
        $("#secretcode .token_url").text(url + "secret/" + result["token"]);
        $("#secretcode").removeClass("invisible");
    })

}

function unlock_secret() {
    var password = $("[name=password]").val();
    var url = window.location.href + "/unlock";
    data = {"password": password}
    send_ajax(data, url, function(result) {
        $("#error").addClass("invisible");
        $("#the-secret p").text(result["secret"]);
        $("#the-secret").removeClass("invisible");
    })   
}