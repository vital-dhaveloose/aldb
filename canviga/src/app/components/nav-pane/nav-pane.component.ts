import { Component } from '@angular/core';
import { fluentTreeView, fluentTreeItem, provideFluentDesignSystem } from "@fluentui/web-components";
import { CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { AldbService } from '../../services/aldb.service';
import { ActivityVersion } from '../../models/activityVersion';
import { CommonModule } from '@angular/common';

provideFluentDesignSystem().register(fluentTreeView(), fluentTreeItem());

@Component({
  selector: 'nav-pane',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './nav-pane.component.html',
  styleUrl: './nav-pane.component.css',
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class NavPaneComponent {

  constructor(
    private aldbService: AldbService
  ) {
  }

  getRootActivity(): ActivityVersion {
    var rootActivity = this.aldbService.getMainActivity();
    console.log(rootActivity)
    return rootActivity
  }

}
