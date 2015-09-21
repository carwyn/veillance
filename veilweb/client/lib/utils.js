

// Based on http://james.padolsey.com/snippets/wordwrap-for-javascript/
// Which in turn is based on http://php.net/manual/en/function.wordwrap.php

function wordwrap(str, width, cut) {
     
    width = width || 75;
    cut = cut || false;
                 
    if (!str) { return str; }
                     
    var regex = '.{1,' + width + '}(\s|$)'
        + (cut ? '|.{' + width + '}|.+$' : '|\S+?(\s|$)');
                         
    return str.match(RegExp(regex, 'g'));                         
}

var debugstring = "TRELAWNEY Dr Livesey, and the rest of these gentlemen, having asked me to write down the whole particulars about Treasure Island, from the beginning to the end, keeping nothing back but the bearings of the island, and that only because there, is still treasure not yet lifted, I take up my pen in the year of, grace 17__ and go back to the time when my father kept the Admiral Benbow inn, and the brown old seaman with the sabre cut first took up his lodging under our roof, I remember him as if it were yesterday, as he came plodding to the inn door, his sea-chest following behind him in a hand-barrow—a tall, strong, heavy, nut-brown man, his tarry pigtail falling over the shoulder, of his soiled blue coat, his hands ragged and scarred, with black, broken nails, and the sabre cut across one cheek, a dirty, livid white. I remember him looking round the cover and whistling to himself as he did so, and then breaking out in that old sea-song that he sang so often afterwards";


//for(var a in (wordwrap(debugstring))) {
//
//   print(a);
//}

