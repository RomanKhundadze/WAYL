function updateText() {
    fetch('/get-playback-data')
        .then(response => response.json())
        .then(data => {
            const trackElement = document.getElementById('track');
            const artistElement = document.getElementById('artist');
            const scrollContainer = document.getElementById('track-container');

            if (data.is_playing) {
                trackElement.textContent = data.item.name;
                artistElement.textContent = data.item.artists.map(artist => artist.name).join(", ");
                adjustScrollDuration();

                const textWidth = trackElement.scrollWidth;
                const containerWidth = scrollContainer.offsetWidth;
                if (textWidth > containerWidth) {
                    scrollContainer.classList.add('scroll');
                } else {
                    scrollContainer.classList.remove('scroll');
                }
            } else {
                trackElement.textContent = "No track currently playing.";
                artistElement.textContent = "";
                scrollContainer.classList.remove('scroll');
            }
        })
        .catch(error => console.error('Error fetching playback state:', error));
}

function adjustScrollDuration() {
    const trackElement = document.getElementById('track');
    const containerWidth = trackElement.parentElement.offsetWidth;
    const textWidth = trackElement.scrollWidth;
    const duration = (textWidth / containerWidth) * 10;

    trackElement.style.animationDuration = `${duration}s`;
}

window.addEventListener('load', updateText);
setInterval(updateText, 1000);
