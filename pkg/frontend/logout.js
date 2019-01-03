function logout() {
    $.ajax({
        type: "GET",
        url: "URL",
        dataType: 'json',
        async: true,
        username: "some_username_noone_has",
        password: "any_password_not_used",
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