import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyDialogModule as MatDialogModule } from '@angular/material/legacy-dialog';
import { TranslateModule } from '@ngx-translate/core';

import { InputModule } from '../input/input.module';
import { NameDialogComponent } from './name-dialog.component';

@NgModule({
  declarations: [NameDialogComponent],
  imports: [CommonModule, MatDialogModule, MatButtonModule, TranslateModule, InputModule, FormsModule],
  exports: [NameDialogComponent],
})
export class NameDialogModule {}
