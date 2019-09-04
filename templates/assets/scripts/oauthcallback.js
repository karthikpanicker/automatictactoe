$( document ).ready(function() {
    if ($('#callback-success').length) {
        window.opener.location.href = "/details";
        self.close();
    }
});