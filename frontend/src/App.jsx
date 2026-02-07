import { useState, useEffect, useCallback } from 'react'
import './App.css'

const HANGMAN_PARTS = [
  // Head
  (props) => <circle cx="70" cy="40" r="20" {...props} />,
  // Body
  (props) => <line x1="70" y1="60" x2="70" y2="130" {...props} />,
  // Left arm
  (props) => <line x1="70" y1="80" x2="40" y2="100" {...props} />,
  // Right arm
  (props) => <line x1="70" y1="80" x2="100" y2="100" {...props} />,
  // Left leg
  (props) => <line x1="70" y1="130" x2="45" y2="170" {...props} />,
  // Right leg
  (props) => <line x1="70" y1="130" x2="95" y2="170" {...props} />,
]

function App() {
  const [word, setWord] = useState('')
  const [guessedLetters, setGuessedLetters] = useState(new Set())
  const [wrongGuesses, setWrongGuesses] = useState(0)
  const [gameStatus, setGameStatus] = useState('loading') // loading, playing, won, lost
  const [difficulty, setDifficulty] = useState('easy')
  const [language, setLanguage] = useState('en')
  const [movieOverview, setMovieOverview] = useState('')
  const [error, setError] = useState('')

  const fetchMovie = useCallback(async () => {
    setGameStatus('loading')
    setError('')
    setGuessedLetters(new Set())
    setWrongGuesses(0)

    try {
      const response = await fetch(`/movie?difficulty=${difficulty}&language=${language}`)
      if (!response.ok) {
        throw new Error('Failed to fetch movie')
      }
      const data = await response.json()
      const title = data.title || ''
      const cleanTitle = title.replace(/[^a-zA-Z\s]/g, '').toUpperCase()
      setWord(cleanTitle)
      setMovieOverview(data.overview || '')
      setGameStatus('playing')
    } catch (err) {
      setError('Failed to load movie. Please try again.')
      setGameStatus('idle')
    }
  }, [difficulty, language])

  useEffect(() => {
    fetchMovie()
  }, [fetchMovie])

  const handleGuess = (letter) => {
    if (gameStatus !== 'playing' || guessedLetters.has(letter)) return

    const newGuessed = new Set(guessedLetters)
    newGuessed.add(letter)
    setGuessedLetters(newGuessed)

    if (!word.includes(letter)) {
      const newWrong = wrongGuesses + 1
      setWrongGuesses(newWrong)
      if (newWrong >= HANGMAN_PARTS.length) {
        setGameStatus('lost')
      }
    } else {
      const letters = word.split('').filter(l => l !== ' ')
      const allGuessed = letters.every(l => newGuessed.has(l))
      if (allGuessed) {
        setGameStatus('won')
      }
    }
  }

  const keyboard = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'.split('')

  return (
    <div className="app">
      <h1>🎬 Movie Hangman</h1>
      
      <div className="controls">
        <label>
          Difficulty:
          <select value={difficulty} onChange={(e) => setDifficulty(e.target.value)}>
            <option value="easy">Easy</option>
            <option value="medium">Medium</option>
            <option value="hard">Hard</option>
          </select>
        </label>
        <label>
          Language:
          <select value={language} onChange={(e) => setLanguage(e.target.value)}>
            <option value="en">English</option>
            <option value="ml">Malayalam</option>
            <option value="hi">Hindi</option>
            <option value="ta">Tamil</option>
          </select>
        </label>
        <button onClick={fetchMovie} className="new-game-btn">New Game</button>
      </div>

      {error && <div className="error">{error}</div>}

      {gameStatus === 'loading' && <div className="loading">Loading movie...</div>}

      {gameStatus !== 'loading' && gameStatus !== 'idle' && (
        <>
          <div className="game-container">
            <div className="hangman-display">
              <svg width="140" height="200" className="hangman-svg">
                {/* Gallows */}
                <line x1="20" y1="190" x2="120" y2="190" stroke="#333" strokeWidth="4" />
                <line x1="70" y1="190" x2="70" y2="20" stroke="#333" strokeWidth="4" />
                <line x1="70" y1="20" x2="70" y2="5" stroke="#333" strokeWidth="4" />
                <line x1="70" y1="5" x2="110" y2="5" stroke="#333" strokeWidth="4" />

                {/* Hangman parts */}
                {HANGMAN_PARTS.slice(0, wrongGuesses).map((Component, index) => (
                  <Component key={index} stroke="#e74c3c" strokeWidth="4" fill="none" />
                ))}
              </svg>
            </div>

            <div className="word-display">
              {word.split('').map((letter, index) => {
                if (letter === ' ') {
                  return <span key={index} className="space"> </span>
                }
                return (
                  <span key={index} className={guessedLetters.has(letter) ? 'revealed' : 'hidden'}>
                    {guessedLetters.has(letter) ? letter : '_'}
                  </span>
                )
              })}
            </div>
          </div>

          {movieOverview && gameStatus === 'won' && (
            <div className="movie-info">
              <h3>Movie Overview:</h3>
              <p>{movieOverview}</p>
            </div>
          )}

          <div className="keyboard">
            {keyboard.map(letter => {
              const isGuessed = guessedLetters.has(letter)
              const isWrong = !word.includes(letter) && isGuessed
              return (
                <button
                  key={letter}
                  onClick={() => handleGuess(letter)}
                  disabled={isGuessed || gameStatus !== 'playing'}
                  className={`key ${isGuessed ? (isWrong ? 'wrong' : 'correct') : ''}`}
                >
                  {letter}
                </button>
              )
            })}
          </div>

          {gameStatus === 'won' && (
            <div className="result won">
              <h2>🎉 Congratulations! You Won!</h2>
              <button onClick={fetchMovie}>Play Again</button>
            </div>
          )}

          {gameStatus === 'lost' && (
            <div className="result lost">
              <h2>😢 Game Over! The word was: {word}</h2>
              <button onClick={fetchMovie}>Try Again</button>
            </div>
          )}

          <div className="stats">
            Wrong guesses: {wrongGuesses} / {HANGMAN_PARTS.length}
          </div>
        </>
      )}
    </div>
  )
}

export default App
