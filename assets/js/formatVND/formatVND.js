
function formatVND(amount) {
    return new Intl.NumberFormat('vi-VN', {
        style:'currency',
        currency:'VND'
    }).format(amount);
}

window.addEventListener('DOMContentLoaded', () => {
    const elements = document.querySelectorAll(".vnd");
    elements.forEach(el => {
        let raw = parseInt(el.innerText.replace(/\D/g, ''));
        if(!isNaN(raw)) {
            el.innerText = formatVND(raw);
        }
    })
})