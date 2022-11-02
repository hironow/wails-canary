import { useEffect, useState, ChangeEvent } from "react";
import "./App.css";
import {
  GetBreedList,
  GetImageUrlsByBreed,
  GetRandomImageUrl,
} from "../wailsjs/go/main/App";

function App() {
  const [randomImageUrl, setRandomImageUrl] = useState("");
  const [breeds, setBreeds] = useState<string[]>([]);
  const [photos, setPhotos] = useState<string[]>([]);

  const [selectedBreed, setSelectedBreed] = useState<string>("");
  const updateSelectedBreed = (event: ChangeEvent<HTMLSelectElement>) => {
    setSelectedBreed(event.target.value);
  };

  const [showRandomPhoto, setShowRandomPhoto] = useState(false);
  const [showBreedPhotos, setShowBreedPhotos] = useState(false);

  useEffect(() => {
    getBreedList();
  }, []);

  function getRandomImageUrl() {
    setShowRandomPhoto(false);
    setShowBreedPhotos(false);
    GetRandomImageUrl().then((result) => setRandomImageUrl(result));
    setShowRandomPhoto(true);
  }

  function getBreedList() {
    GetBreedList().then((result) => setBreeds(result));
  }

  function getImageUrlsByBreed() {
    setShowRandomPhoto(false);
    setShowBreedPhotos(false);
    GetImageUrlsByBreed(selectedBreed).then((result) => setPhotos(result));
    setShowBreedPhotos(true);
  }

  return (
    <div id="App">
      <div>
        <button className="btn" onClick={getRandomImageUrl}>
          Fetch a dog randomly
        </button>
        <p>Click on down arrow to select a breed</p>
        <select
          title="select-breed"
          value={selectedBreed}
          onChange={updateSelectedBreed}
        >
          {breeds?.map((breed, i) => (
            <option key={i} value={breed}>
              {breed}
            </option>
          ))}
        </select>
        <button className="btn" onClick={getImageUrlsByBreed}>
          Fetch by this breed
        </button>
      </div>

      {showRandomPhoto && (
        <img id="random-photo" src={randomImageUrl} alt="No dog found" />
      )}

      {showBreedPhotos &&
        photos?.map((photo, i) => (
          <img key={i} id="breed-photos" src={photo} alt="No dog found" />
        ))}
    </div>
  );
}

export default App;
