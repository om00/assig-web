function toggleOtherReasonInput() {
    var selectedReason = document.getElementById('blockReasonSelect').value;
    var otherReasonContainer = document.getElementById('otherReasonContainer');
    if (selectedReason == '99') {
        otherReasonContainer.style.display = 'block';
    } else {
        otherReasonContainer.style.display = 'none';
    }
}

// On Modal Show, Pass User ID to Modal
$('#blockModal').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget); // Button that triggered the modal
    var userId = button.data('user-id'); // Extract info from data-* attributes
    // Optionally, store the userId in a hidden field or use it in the request
    $('#confirmBlockBtn').data('user-id', userId);
});

// Handle Confirm Block Button
$('#confirmBlockBtn').on('click', function() {
    var userId = $(this).data('user-id');
    var selectedReason = $('#blockReasonSelect').val();
    var otherReason = $('#otherReasonInput').val();
    var reason
    var data

    if (!selectedReason) { 
        alert('Please select any reason');
        return;
    }

    
    if (selectedReason === '99' && otherReason === '') {
        alert('Please provide a reason if "Other" is selected.');
        return;
    }
    
    if (selectedReason === '99'){
        reason = otherReason
    }
    
    if (!userId){
         data = {
            name: $('#username').val(),
            phone: $('#phone').val(),
            email: $('#email').val(),
            status: $('#status').val(),
            blockReason: $('#blockReason').val(),
            reasonCode:  selectedReason,
            reason:reason
        };
        $('#blockAllUsersBtn').text('Processing...').prop('disabled', true);

    }else{
        data = {
            userId: userId,
            reasonCode:  selectedReason,
            reason:reason
        }

    }
  

    // Example: Send request to server with selected reason
    $.ajax({
        url: '/block-user', // Server endpoint to handle block
        method: 'POST',
        
        data: JSON.stringify(data),
        
        success: function(response) {
            
            Swal.fire({
                title: 'Success!',
                text: response.message,
                icon: 'success',
                timer: 3000,  // Auto close after 2 seconds
                showConfirmButton: false
            });
            
            if (!userId){
                $('#blockAllUsersBtn').text('Block Selected  Users').prop('disabled', false);
            }
            setTimeout(function() {
                // Hide the message box
                location.reload();  // Reload the page
            }, 4000);
        },
        error: function() {
            
            Swal.fire({
                title: 'Error!',
                text: 'An error occurred while blocking the user.',
                icon: 'error',
                timer: 2000,  // Auto close after 2 seconds
                showConfirmButton: false
            });
            if (!userId){
                $('#blockAllUsersBtn').text('Block Selected  Users').prop('disabled', false);
            }
           
        }
    });
});

