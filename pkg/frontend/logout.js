function logout() {
    $.ajax({
        type: "GET",
        url: "PUT_YOUR_PROTECTED_URL_HERE",
        dataType: 'json',
        async: true,
        username: "some_username_that_doesn't_exist",
        password: "any_stupid_password",
        data: '{ "comment" }'
    })
        .done(function(){
            alert('Der genutzte Browser ist zu alt oder zu speziell und unterstützt den Logout nicht. Bitte schließen sie ihren Browser oder löschen sie ihren Authentifizierungs-Cache.')
        })
        .fail(function(){
            alert('Logout abgeschlossen.');
            window.location = "/";
        });
}