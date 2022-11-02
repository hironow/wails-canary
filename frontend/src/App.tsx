import { useEffect, useState, ChangeEvent } from 'react'
import './App.css'
import { GetBreedList, GetImageUrlsByBreed, GetRandomImageUrl } from '../wailsjs/go/main/App'

function App() {
  // state
  const [randomImageUrl, setRandomImageUrl] = useState('')
  const [breeds, setBreeds] = useState<string[]>([])
  const [photos, setPhotos] = useState<string[]>([])

  // select
  const [selectedBreed, setSelectedBreed] = useState<string>('')
  const updateSelectedBreed = (event: ChangeEvent<HTMLSelectElement>) => {
    setSelectedBreed(event.target.value)
  }

  // flags
  const [showFlags, setShowFlags] = useState<{
    randomPhoto: boolean
    breedPhoto: boolean
  }>({
    randomPhoto: false,
    breedPhoto: false,
  })

  // init
  useEffect(() => {
    getBreedList()
  }, [])

  // call
  function getRandomImageUrl() {
    setShowFlags({ randomPhoto: false, breedPhoto: false })
    GetRandomImageUrl().then((result) => setRandomImageUrl(result))
    setShowFlags({ randomPhoto: true, breedPhoto: false })
  }

  function getBreedList() {
    GetBreedList().then((result) => setBreeds(result))
  }

  function getImageUrlsByBreed() {
    setShowFlags({ randomPhoto: false, breedPhoto: false })
    GetImageUrlsByBreed(selectedBreed).then((result) => setPhotos(result))
    setShowFlags({ randomPhoto: false, breedPhoto: true })
  }

  return (
    <div id='App'>
      <div id='input' className='input-box'>
        <button type='button' className='btn' onClick={getRandomImageUrl}>
          Fetch a dog randomly
        </button>
        <p>Click on down arrow to select a breed</p>
        <select title='select-breed' value={selectedBreed} onChange={updateSelectedBreed}>
          {breeds?.map((breed, i) => (
            <option key={i} value={breed}>
              {breed}
            </option>
          ))}
        </select>
        <button type='button' className='btn' onClick={getImageUrlsByBreed}>
          Fetch by this breed
        </button>
      </div>

      {showFlags.randomPhoto && <img id='random-photo' src={randomImageUrl} alt='No dog found' />}

      {showFlags.breedPhoto &&
        photos?.map((photo, i) => <img key={i} id='breed-photos' src={photo} alt='No dog found' />)}
    </div>
  )
}

export default App
