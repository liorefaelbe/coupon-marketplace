import { type FormEvent, useState } from "react";

type AdminLoginProps = {
  onLogin: (username: string, password: string) => Promise<void>;
  onCancel: () => void;
};

function AdminLogin({ onLogin, onCancel }: AdminLoginProps) {
  const [username, setUsername] = useState("admin");
  const [password, setPassword] = useState("admin123");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      setLoading(true);
      setError("");
      await onLogin(username, password);
    } catch (err) {
      const message = err instanceof Error ? err.message : "Login failed";
      setError(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modalOverlay">
      <div className="modalCard">
        <h2>Admin Login</h2>
        <p>Enter admin credentials to access the management panel.</p>

        <form className="form" onSubmit={handleSubmit}>
          <input placeholder="Admin username" value={username} onChange={(e) => setUsername(e.target.value)} required />
          <input
            type="password"
            placeholder="Admin password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />

          <div className="modalActions">
            <button type="button" className="secondaryButton" onClick={onCancel}>
              Cancel
            </button>
            <button type="submit" disabled={loading}>
              {loading ? "Checking..." : "Login"}
            </button>
          </div>
        </form>

        {error && <div className="message error">{error}</div>}
      </div>
    </div>
  );
}

export default AdminLogin;
