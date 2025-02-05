const params = new URLSearchParams(window.location.search);
const tourId = params.get("id");
fetch(`/api/tours/${tourId}`)
    .then(response => response.json())
    .then(tour => {
        document.getElementById("tour-details").innerHTML = `
            <h2>${tour.name}</h2>
            <p>${tour.description}</p>
        `;
    });