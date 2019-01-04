function logout() {
    try {
        // for Firefox
        $.ajax({
            url: "/secure/dashboard.html",
            username: 'please enter your username here',
            password: 'any_password_not_used',

            statusCode: {
                401: function () {
                    document.location = "/";
                }
            }

        });

    } catch (exception) {
        // for IE
        if (!document.execCommand("ClearAuthenticationCache")) {
            //for Chrome
            document.location = "https://reset:reset@" + document.location.hostname + document.location.pathname;
        }
    }

    document.location = "/";
}
