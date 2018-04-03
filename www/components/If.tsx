import * as React from 'react';

/**
 * Usage:
 * <If condition={true}>
 *  <span>I will render.</span>
 * </If>
 *
 */
export function If(props: { condition: boolean, children?: any }) {
  return props.condition ? props.children : null;
}
