function checkThemeCookie(cookie) {
    var val = Cookies.get(cookie);
    if(typeof val == 'string' && val.length > 0) {
        changeTheme(val, cookie);
    }
}

function changeTheme(theme, cookie) {
    var body = $('body');
    body.removeClass(body.attr('data-theme'));
    body.addClass(theme).attr('data-theme', theme);

    if(typeof cookie == 'string' && cookie.length > 0) {
        Cookies.set(cookie, theme, {expires: Infinity});
    }
}

function newLiveClient(handler) {
    var socket = io('http://realmstatus.net');
    socket.on('realms', function (data) {
        handler(data)
    });
}

$(document).ready(function() {
    var cookie = 'realmstatus_net_theme';

    checkThemeCookie(cookie);

    $('.theme-changer').click(function(e) {
        changeTheme($(this).attr('data-theme'), cookie);
        e.preventDefault();
    });
});
