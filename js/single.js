$(document).ready(function() {
    // Start a timer to continually display and update the realm time
    setInterval(function() {
        var now = moment();
        var tz = $('#realm-timezone').text();
        var timeString = now.tz(tz).format('ddd Do MMM YYYY h:m:s A');
        $('#realm-time').text(timeString);
        $('#realm-timezone-short').text(now.tz(tz).format('(z)'));
    }, 1000);

    // Bind the button to show/hide battlegroup realms
    $('a.battlegroup-toggler').click(function(e) {
        var list = $('ul.battlegroup-list');
        var icon = '+';
        list.toggleClass('visible');
        if(list.hasClass('visible')) {
            icon = '-';
        }
        $(this).children('span.icon').text(icon)
        e.preventDefault();
    });

    // Start live data client to get live updates
    newLiveClient(function(data) {
        for(var realm in data) {
            if(realm.hasOwnProperty(property)) {
                if(realm.Status) {
                    $('.live-status').html('<span class="good">UP</span>');
                }
                else {
                    $('.live-status').html('<span class="bad">DOWN</span>');
                }

                if(realm.Queue) {
                    $('.live-queue').html('<span class="bad">NO</span>');
                }
                else {
                    $('.live-queue').html('<span class="good">YES</span>');
                }
            }
        }
    });
});