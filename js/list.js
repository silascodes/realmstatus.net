// Function to highlight odd visible table rows
function highlightList() {
    $('#list table tbody tr.visible').each(function(i, e) {
        if(i % 2 == 0) {
            $(e).addClass("highlight");
        }
        else {
            $(e).removeClass("highlight");
        }
    });
}

// Function to apply selected filters to list table
function applyFilters() {
    // Get values of filters on the page that have non-default and non-blank values
    var filters = []
    $('#filters .filter').each(function(i, e) {
        var v = $(e).val();
        if(v.length > 0 && v != 'any') {
            filters.push({
                'field': $(e).attr('data-field'),
                'value': v,
            });
        }
    });

    // Apply the filters to the table
    if(filters.length > 0) {
        // Check each table row
        $('#list table tbody tr').each(function(i, e) {
            // Find out how many filters match this row
            var lineMatches = 0;
            filters.forEach(function(pe, pi, pa) {
                var v = $(e).children('td[data-field="' + pe.field + '"]').attr('data-value');
                if(v == pe.value) {
                    lineMatches++;
                }
            });

            // If all the filters match, make it visible, otherwise make it hidden
            if(lineMatches < filters.length) {
                $(e).removeClass('visible');
            }
            else {
                $(e).addClass('visible');
            }
        });

        // Run the table row highlighter
        highlightList();
    }
}

$(document).ready(function() {
    // Do initial table row highlighting
    highlightList();

    // Initialise table sorter, sort by status (asc) then name (asc)
    $('#list table').tablesorter({
        sortList: [[6,0], [1,0]]
    });

    // Hook table sortEnd event to call highlighter
    $('#list table').on('sortEnd', function(e) {
        highlightList();
    });

    // Hook apply filters button (hide/show table rows based on filter values)
    $('#apply-filters').on('click', function(e) {
        applyFilters();
        e.preventDefault();
    });

    // Hook reset filters button (reload the page)
    $('#reset-filters').on('click', function(e) {
        location.reload();
        e.preventDefault();
    });

    // Start live data client to get live updates
    newLiveClient(function(data) {
        for(var realm in data) {
            if(realm.hasOwnProperty(property)) {
                var statusContent = '<span class="bad">DOWN</span>';
                if(realm.Status) {
                    statusContent = '<span class="good">UP</span>';
                }

                var queueContent = '<span class="good">YES</span>';
                if(realm.Queue) {
                    queueContent = '<span class="bad">NO</span>';
                }

                var tr = $('tr[data-slug="' + i + '"]');
                tr.children('td.live-status').html(statusContent);
                tr.children('td.live-queue').html(queueContent);
            }
        }
    });
});