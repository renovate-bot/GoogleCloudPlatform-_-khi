/**
 * Copyright 2024 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { NgModule } from '@angular/core';
import {
  LongTimestampFormatPipe,
  TimestampFormatPipe,
} from './timestamp-format.pipe';
import { FirstOrUndefined } from './first-or-null.pipe';
import { ParsePathPipe } from './parse-path.pipe';
import { RainbowPipe } from './rainbow.pipe';
import { SidePaneComponent } from './components/side-pane.component';
import { ResolveTextPipe } from './resolve-text.pipe';
import { MetaTableRowComponent } from './components/meta-table-row.component';
import { CssClassFormatPipe } from './css-class-format.pipe';
import { BreaklinePipe } from './breakline.pipe';
import { CaptureShiftKeyDirective } from './capture-shiftkey.directive';

@NgModule({
  imports: [
    TimestampFormatPipe,
    LongTimestampFormatPipe,
    FirstOrUndefined,
    ParsePathPipe,
    RainbowPipe,
    ResolveTextPipe,
    CssClassFormatPipe,
    SidePaneComponent,
    MetaTableRowComponent,
    BreaklinePipe,
    CaptureShiftKeyDirective,
  ],
  exports: [
    TimestampFormatPipe,
    LongTimestampFormatPipe,
    FirstOrUndefined,
    ParsePathPipe,
    RainbowPipe,
    ResolveTextPipe,
    CssClassFormatPipe,
    SidePaneComponent,
    MetaTableRowComponent,
    BreaklinePipe,
    CaptureShiftKeyDirective,
  ],
})
export class KHICommonModule {}
