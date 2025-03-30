const API_URL = 'http://localhost:8082'; // Или адрес бекенда в Docker-сети

const cityInput = document.getElementById('cityInput');
const suggestionsDiv = document.getElementById('suggestions');
const resultDiv = document.getElementById('result');
const recommendationText = document.getElementById('recommendationText');

let timeoutId;

// Поиск городов с задержкой 300 мс
cityInput.addEventListener('input', (e) => {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => searchCities(e.target.value), 300);
});

// Выбор города из списка
suggestionsDiv.addEventListener('click', (e) => {
    if (e.target.classList.contains('suggestion-item')) {
        const location = JSON.parse(e.target.dataset.location);
        console.log('Выбранный город:', location);
        getRecommendation(location);
    }
});

async function searchCities(query) {
    if (!query) {
        suggestionsDiv.style.display = 'none';
        return;
    }

    try {
        const response = await fetch(`${API_URL}/geocode?query=${encodeURIComponent(query)}`);
        const data = await response.json();
        console.log('Ответ от /geocode:', data);

        suggestionsDiv.innerHTML = data.results
            .map(loc => `
                <div class="suggestion-item" 
                     data-location='${JSON.stringify(loc)}'>
                    ${loc.name}
                </div>
            `).join('');

        suggestionsDiv.style.display = 'block';
    } catch (error) {
        console.error('Ошибка поиска:', error);
    }
}

async function getRecommendation(location) {
    cityInput.value = location.name;
    suggestionsDiv.style.display = 'none';
    resultDiv.classList.add('hidden');

    try {
        const response = await fetch(`${API_URL}/recommend`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({
                latitude: location.latitude,
                longitude: location.longitude
            })
        });

        const data = await response.json();
        recommendationText.textContent = data.data;
        resultDiv.classList.remove('hidden');
    } catch (error) {
        recommendationText.textContent = 'Ошибка получения рекомендаций';
        resultDiv.classList.remove('hidden');
    }
}