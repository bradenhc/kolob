//------------------------------------------------------------------------------------------------//
//--  Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)   --//
//------------------------------------------------------------------------------------------------//
import { LitElement, css, html } from "lit";

export class KolobSignIn extends LitElement {
  static properties = {
    _group: { state: true },
    _username: { state: true },
    _password: { state: true },
    _errors: { state: true },
  };

  constructor() {
    super();
    this._group = "";
    this._username = "";
    this._password = "";
    this._errors = [];
  }

  render() {
    return html`
      <h2>Welcome to Kolob!</h2>
      <form @submit=${this._handleSubmit}>
        <label for="group">Group</label>
        <input type="text" id="group" @change=${this._handleGroupChange} />
        <label for="username">Username</label>
        <input type="text" id="username" @change=${this._handleUsernameChange} />
        <label for="password">Password</label>
        <input type="password" id="password" @change=${this._handlePasswordChange} />
        <input type="submit" hidden />
        <button type="submit">Sign In</button>
      </form>
      ${this._renderErrors()}
    `;
  }

  _renderErrors() {
    if (this._errors.length === 0) {
      return "";
    }

    return html`
      <div class="error-container">
        ${this._errors.map((e) => html`<div class="error">${e}</div>`)}
      </div>
    `;
  }

  _handleGroupChange(e) {
    this._group = e.target.value;
  }

  _handleUsernameChange(e) {
    this._username = e.target.value;
  }

  _handlePasswordChange(e) {
    this._password = e.target.value;
  }

  _handleSubmit(e) {
    e.preventDefault();
    const errors = [];
    if (this._group === "") {
      errors.push("Missing group name");
    }
    if (this._username === "") {
      errors.push("Missing username");
    }
    if (this._password === "") {
      errors.push("Missing password");
    }
    if (errors.length !== 0) {
      this._errors = errors;
      return false;
    }

    this._errors = [];
    const detail = { username: this._username, name: "Braden Hitchcock" };
    const event = new CustomEvent("kolob-authenticated", { bubbles: true, detail });
    this.dispatchEvent(event);
    return false;
  }
}

KolobSignIn.styles = css`
  :host {
    align-self: center;
    justify-self: center;
  }

  form {
    display: flex;
    flex-direction: column;
  }

  label {
    vertical-align: top;
  }

  input {
    margin: 5px;
  }

  button {
    align-self: center;
    width: 50%;
    margin-top: 12px;
  }

  .error-container {
    display: flex;
    flex-direction: column;
  }

  .error {
    color: red;
  }
`;
