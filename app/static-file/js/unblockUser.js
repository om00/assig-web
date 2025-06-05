$(document).ready(function () {
    // When the Unblock button is clicked
    $('.unblockButton').click(function () {
        // Prepare the data to send (you can use dynamic data here)
        var userId = $(this).data('user-id');
        var data = {
            userId: userId
        };

        
        $.ajax({
            url: '/unblock-user', 
            method: 'POST',  
            contentType: 'application/json',  
            data: JSON.stringify(data),  
            success: function (response) {
                // On success, show a success message and reload the page (if needed)
                Swal.fire({
                    title: 'Success!',
                    text: 'User unblocked successfully!',
                    icon: 'success',
                    timer: 2000,  // Auto close after 2 seconds
                    showConfirmButton: false
                });

                // Optionally, reload the page after a short delay
                setTimeout(function () {
                    location.reload();
                }, 3000);
            },
            error: function (xhr, status, error) {
                // On error, show an error message
                Swal.fire({
                    title: 'Error!',
                    text: 'An error occurred while unblocking the user.',
                    icon: 'error',
                    timer: 2000,  // Auto close after 2 seconds
                    showConfirmButton: false
                });

            }
        });
    });

        const $modal = $('#confirmModal');
      
        $('#unblockAllUsersBtn').click(function() {
          $modal.show(); 
        });
      
        $('#confirmNo').click(function() {
          $modal.hide(); 
        });
      
        $('#confirmYes').click(function() {
          $modal.hide(); 

          data = {
            name: $('#username').val(),
            phone: $('#phone').val(),
            email: $('#email').val(),
            status: $('#status').val(),
            blockReason: $('#blockReason').val(),
        };
      
          // 🚀 jQuery AJAX
          $.ajax({
            url: '/unblock-user', 
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify( data),
            success: function(response) {
             
             Swal.fire({
                title: 'Success!',
                text: 'User unblocked successfully!',
                icon: 'success',
                timer: 2000,  // Auto close after 2 seconds
                showConfirmButton: false
            });

            // Optionally, reload the page after a short delay
            setTimeout(function () {
                location.reload();
            }, 3000);
            },
            error: function(xhr, status, error) {
               // On error, show an error message
               Swal.fire({
                title: 'Error!',
                text: 'An error occurred while unblocking the user.',
                icon: 'error',
                timer: 2000,  // Auto close after 2 seconds
                showConfirmButton: false
            });
            }
          });
        });
      
});