import { BrowserRouter, Link, Route, Routes } from "react-router-dom";

function HomePage() {
  return (
    <main>
      <h1>Amara</h1>
      <p>The application workspace is ready for feature development.</p>
      <Link to="/jobs">Browse jobs</Link>
    </main>
  );
}

function JobsPage() {
  return (
    <main>
      <h1>Jobs</h1>
      <Link to="/">Back home</Link>
    </main>
  );
}

export function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/jobs" element={<JobsPage />} />
      </Routes>
    </BrowserRouter>
  );
}