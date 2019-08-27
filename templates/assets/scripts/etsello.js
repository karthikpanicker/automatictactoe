 $( document ).ready(function() {
    $('#boards button').on("click",function(){ 
        var boardId = $(this).attr("id");
        loadBoardLists(boardId,$(this));
     });
     if ($('#boards button').length > 0) {
        var firstElem = $('#boards button')[0];
        var boardId = $(firstElem).attr('id');
        loadBoardLists(boardId,$(firstElem));
     }

     $("#saveList").on("click",function(){
        var $this = $(this);
        var loadingText = '<i class="fa fa-circle-o-notch fa-spin"></i> Saving...';
        if ($(this).html() !== loadingText) {
            $this.data('original-text', $(this).html());
            $this.html(loadingText);
        }
        var bId = $('#boards button').filter('.active').attr('id');
        var lId = $('#boardLists button').filter('.active').attr('id');
        $.ajax({
            type: "POST",
            url: "api/user-info",
            data: JSON.stringify({"boardId": bId, "listId": lId}),
            success: function(data){
                setTimeout(function () {
                    $this.html($this.data('original-text'));
                }, 1000);
            }
          });
    });
});

function loadBoardLists(boardId,selectedBoard) {
    $('#spinner-board-list').show();
    $('#boards button').removeClass('active');
    selectedBoard.addClass('active');
    $.get( "api/trello-boards/" + boardId + "/lists", function( boardLists ) {
        $('#boardLists button').remove();
        $.each(boardLists, function( index, list ) {
            $('#boardLists').append('<button id="' + list.id 
            + '" type="button" class="list-group-item">' + list.name + '</button>');
        });
        if ($('#boardLists button').length > 0) {
            $('#boardLists button').eq(0).addClass('active');
            // Bind click action on newly added items
            $('#boardLists button').on("click",function(){ 
                $('#boardLists button').removeClass('active');
                $(this).addClass('active');
                $('#spinner-board-list').hide();
            });
        }
      }); 
}