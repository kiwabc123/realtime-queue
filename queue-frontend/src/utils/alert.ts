import Swal from 'sweetalert2'

export function showAlert(msg: string) {
  Swal.fire({
    text: msg,
    icon: 'info',
    timer: 1800,
    showConfirmButton: false,
    position: 'top',
    toast: true,
    timerProgressBar: true
  })
}
