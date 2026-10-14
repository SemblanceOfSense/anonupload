async function loadImages() {
    const div = document.getElementById("displayedImagesDiv");
    try {
        const res = await fetch("/listfiles");
        if (!res.ok) throw new Error(await res.text());
        const data = await res.json();

        div.replaceChildren();
        for (const name of data.Files ?? []) {
            const img = document.createElement("img");
            img.src = "/getfile?filename=" + encodeURIComponent(name);
            img.alt = name;
            img.style.maxWidth = "300px";
            img.style.margin = "4px";
            div.appendChild(img);
        }
    } catch (err) {
        div.textContent = "Couldn't load images: " + err.message;
    }
}

document.getElementById("uploadForm").addEventListener("submit", async (e) => {
    e.preventDefault();
    const file = document.getElementById("myFile").files[0];
    if (!file) return;
    const res = await fetch("/upload", { method: "POST", body: file });
    if (!res.ok) alert(await res.text());
    e.target.reset();
    loadImages();
});

loadImages();
