document.addEventListener("DOMContentLoaded", () => {
    fetch("/api/tours")
        .then(response => response.json())
        .then(data => {
            const toursDiv = document.getElementById("tours");
            toursDiv.innerHTML = data.map(tour => `
                <div class='tour'>
                    <h3>${tour.name}</h3>
                    <p>${tour.description}</p>
                    <a href='/tour.html?id=${tour.id}'>View Details</a>
                </div>
            `).join('');
        });
});