//------------------------------------------------------------------------------------------------//
//--  Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)   --//
//------------------------------------------------------------------------------------------------//

/**
 * Defines a new custom element with the registry. If a custom element with the same name already
 * exists, then a new element is not registered. This avoids the exception that is innevitably
 * thrown when you try to register the exact same custom element twice.
 *
 * @param {string}                      name    The name of the custom element.
 * @param {CustomElementConstructor}    ctor    The constructor used to create the custom element.
 */
export function define(name, ctor) {
  if (customElements.get(name) === undefined) {
    customElements.define(name, ctor);
  }
}
