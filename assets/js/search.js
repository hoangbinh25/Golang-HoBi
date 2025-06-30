document.addEventListener('DOMContentLoaded', function () {
    const searchInput = document.getElementById('search-input');
    const suggestionsContainer = document.getElementById('search-suggestions');
    const searchForm = document.getElementById('search-form');

    let debounceTimer;

    searchInput.addEventListener('input', function () {
        const query = this.value.trim();

        clearTimeout(debounceTimer);

        if (query.length < 1) {
            suggestionsContainer.innerHTML = '';
            suggestionsContainer.style.display = 'none';
            return;
        }

        debounceTimer = setTimeout(() => {
            fetch(`/api/search-suggestions?q=${encodeURIComponent(query)}`)
                .then(response => response.json())
                .then(data => {
                    suggestionsContainer.innerHTML = '';
                    if (data && data.length > 0) {

                        const shopSuggestion = `<a href="/search?q=${encodeURIComponent(query)}" class="list-group-item list-group-item-action"><i class="fa-solid fa-store me-2"></i>Tìm Shop "${query}"</a>`;
                        suggestionsContainer.innerHTML += shopSuggestion;

                        data.forEach(suggestion => {
                            const item = document.createElement('a');
                            item.href = `/search?q=${encodeURIComponent(suggestion)}`;
                            item.classList.add('list-group-item', 'list-group-item-action');
                            item.textContent = suggestion;
                            suggestionsContainer.appendChild(item);
                        });
                        suggestionsContainer.style.display = 'block';
                    } else {
                        suggestionsContainer.style.display = 'none';
                    }
                })
                .catch(error => {
                    console.error('Error fetching search suggestions:', error);
                    suggestionsContainer.style.display = 'none';
                });
        }, 300);
    });

    document.addEventListener('click', function (e) {
        if (!searchInput.contains(e.target)) {
            suggestionsContainer.style.display = 'none';
        }
    });

    suggestionsContainer.addEventListener('click', function (e) {
        e.stopPropagation();
    });
}); 