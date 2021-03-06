//2057008, 2624395, 9111696
function logout() {
    //Ausloggfunktion für die Basic Authetication, die mindestens die Browser Firefox, Chrome und Internet Explorer abdeckt (jeweils Versionsunabhängig)
    try {
        // Ausloggen aus der Basic Authentication in Firefox durch Eingabe falscher Daten,
        $.ajax({
            url: "/secure/dashboard.html",
            username: 'bitte Benutzernamen eingeben',
            password: 'irgendein_ungenutztes_Passwort',
        // Code 401 heißt, dass das Ausloggen erfolgreich war - die Exception wird dennoch geworfen (ist so vorgesehen)
            statusCode: {
                401: function () {
                    document.location = "/";
                }
            }
        });
    } catch (exception) {
        // leeren des Caches für den Logout im Internet Explorer
        if (!document.execCommand("ClearAuthenticationCache")) {
            // ersetzen der gespeicherten Basic Authentication-Kennung in Chrome
            document.location = "https://reset:reset@" + document.location.hostname + document.location.pathname;
        }
    }
    document.location = "/";
}
