function logout() {
    try {
        $.ajax({
            url: "/secure/dashboard.html",
            username: 'please enter your username here',
            password: 'any_password_not_used',

            statusCode: {
                401: function () {
                    window.location = "/";
                }
            }

                .done(function(){
                    alert('Der genutzte Browser ist zu alt oder zu speziell und unterstützt den Logout nicht. Bitte schließen sie ihren Browser oder löschen sie ihren Authentifizierungs-Cache.')
                })
                .fail(function(){
                    alert('Logout abgeschlossen.');
                    window.location = "/";
                })

        });
    } catch (exception) {
        if (!document.execCommand("ClearAuthenticationCache")) {
            document.location = "https://reset:reset@" + document.location.hostname + document.location.pathname;
        }
    };

    window.location = "/";
}
