
Meteor.utils = {

    // Based on http://james.padolsey.com/snippets/wordwrap-for-javascript/
    // Which in turn is based on http://php.net/manual/en/function.wordwrap.php

    wordwrap: function (str, width, cut) {
     
        width = width || 50;
        cut = cut || false;
                     
        if (!str) { return str; }
                         
        var regex = '.{1,' + width + '}(\s|$)'
            + (cut ? '|.{' + width + '}|.+$' : '|\S+?(\s|$)');
                             
        return str.match(RegExp(regex, 'g'));                         
    }
}

