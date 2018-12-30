function logout (c) {
    var a,b="Logout is complete";
    try {
        a=document.execCommand("ClearAuthenticationCache")
    }
    catch(d) {
    }
    a||((a=window.XMLHttpRequest?new window.XMLHttpRequest:window.ActiveXObject?new ActiveXObject("Microsoft.XMLHTTP"):void 0)?(a.open("HEAD",c||location.href,!0,"logout",(new Date).getTime().toString()),a.send(""),a=1):a=void 0);
    a||(b="Your browser is too old or too weird to support log out functionality. Close all windows and restart the browser or manually clear your authentication cache.");
    alert(b);
    location.replace("/index.html")
}