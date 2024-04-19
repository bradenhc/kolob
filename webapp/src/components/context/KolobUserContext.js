//------------------------------------------------------------------------------------------------//
//--  Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)   --//
//------------------------------------------------------------------------------------------------//
import { ContextConsumer, ContextProvider, createContext } from "@lit/context";

const KOLOB_USER_CONTEXT = createContext(Symbol("kolob-user-context"));

export class KolobUserContextProvider extends ContextProvider {
  constructor(host) {
    super(host, { context: KOLOB_USER_CONTEXT });
  }
}

export class KolobUserContextConsumer extends ContextConsumer {
  constructor(host) {
    super(host, { context: KOLOB_USER_CONTEXT });
  }
}
